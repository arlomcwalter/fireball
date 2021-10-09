package minegame159.fireball.types;

public class PrimitiveType extends Type {
    public final PrimitiveTypes type;

    public PrimitiveType(String name, PrimitiveTypes type) {
        super(name);
        this.type = type;
    }

    @Override
    public boolean canBeAssignedTo(Type to) {
        if (!(to instanceof PrimitiveType)) return false;
        PrimitiveTypes toType = ((PrimitiveType) to).type;

        return toType.size >= type.size;
    }

    @Override
    protected Type copy() {
        return new PrimitiveType(name, type);
    }

    @Override
    public boolean equals(Type type) {
        if (!super.equals(type)) return false;
        return this.type == ((PrimitiveType) type).type;
    }
}
